import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:math';

class CustomButton extends StatelessWidget {
  const CustomButton({
    super.key,
    required this.text,
    required this.onPressed,
    this.backgroundColor = Colors.white,
    this.size = const Size(double.infinity, 70),
    this.padding = const EdgeInsets.symmetric(),
    this.textPaddingLeft = 0,
    this.image,
    this.borderCircularRadius = 32.0,
  });

  final Color backgroundColor;
  final Widget? text;
  final VoidCallback? onPressed;
  final Size size;
  final EdgeInsets padding;
  final double textPaddingLeft;
  final Image? image;
  final double borderCircularRadius;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: padding,
      child: ElevatedButton(
        style: ElevatedButton.styleFrom(
          backgroundColor: backgroundColor,
          elevation: 3,
          minimumSize: size,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(borderCircularRadius),
          ),
        ),
        onPressed: onPressed,
        child: Stack(
          alignment: Alignment.center,
          children: [
            if (image != null)
              Positioned(
                left: 0,
                child: image!,
              ),
            Padding(padding: EdgeInsets.only(left: textPaddingLeft), child: text),
            const Center(),
          ],
        ),
      ),
    );
  }
}


class CustomClassicForm extends StatelessWidget {
  const CustomClassicForm({
    super.key,
    required this.hint,
    this.controller,
    this.obscureText = false,
    this.padding = const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
    this.enabled = true,
    this.fillCollor = Colors.white,
  });

  final String hint;
  final TextEditingController? controller;
  final bool obscureText;
  final EdgeInsetsGeometry padding;
  final bool enabled;
  final Color fillCollor;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: <Widget>[
        Padding(
          padding: padding,
          child:
          Container(
            decoration: BoxDecoration(
              borderRadius:  BorderRadius.circular(10),
            ),
            child: TextField(
              style: const TextStyle(fontSize:25, fontWeight: FontWeight.w900, color: Colors.black),
              controller: controller,
              obscureText: obscureText,
              cursorColor: Colors.black,
              decoration: InputDecoration(
                filled: true,
                fillColor: fillCollor,
                labelStyle:const TextStyle(fontSize:40),
                hintStyle: const TextStyle(
                  fontSize: 20,
                  color: Color.fromARGB(255, 150, 150, 150),
                  fontWeight: FontWeight.bold,
                ),
                hintText: hint,
                enabled: enabled,
                contentPadding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
                enabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: const BorderSide(
                    color: Color.fromARGB(255, 238, 238, 238),
                    width: 5,
                  ),
                ),
                disabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: const BorderSide(
                    color: Color.fromARGB(255, 238, 238, 238),
                    width: 5,
                  ),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(10),
                  borderSide: const BorderSide(
                    color: Colors.black,
                    width: 5,
                  ),
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }
}

class MyDivider extends StatelessWidget {
  const MyDivider({
    super.key,
    this.text = const AutoSizeText(''),
    this.textPadding = const EdgeInsets.symmetric(),
    this.color = Colors.grey,
    this.thickness = 1,
  });

  final AutoSizeText text;
  final EdgeInsets textPadding;
  final Color color;
  final double thickness;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(child: Divider(thickness: thickness, color: color,),),
        Padding(padding: textPadding, child: text,),
        Expanded(child: Divider(thickness: thickness, color: color,),),
      ],
    );
  }
}

class MyVerticalDivider extends StatelessWidget {
  const MyVerticalDivider({
    super.key,
    this.height,
    this.width,
    this.color = Colors.black,
    this.padding = const EdgeInsets.only(left: 0),
  });

  final double? height;
  final double? width;
  final Color color;
  final EdgeInsetsGeometry padding;

  @override
  Widget build(BuildContext context) {
    return Row (
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        Padding(
          padding: padding,
          child: SizedBox(
            height: height,
            width: width,
            child: DecoratedBox(
              decoration: BoxDecoration(
                color: color,
              ),
            ),
          ),
        ),
      ],
    );
  }
}

class PlaceableWidget extends StatelessWidget {
  const PlaceableWidget({
    super.key,
    required this.child,
    this.left = 0,
    this.right = 0,
    this.top = 0,
    this.bottom = 0,
  });

  final Widget child;
  final double left;
  final double right;
  final double top;
  final double bottom;

  @override
  Widget build(BuildContext context) {
    return Row (
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        Padding(
          padding: EdgeInsets.only(left: left, right: right, top: top, bottom: bottom),
          child: child,
        ),
      ],
    );
  }
}

class LeftWidget extends StatelessWidget {
  const LeftWidget({
    super.key,
    required this.child,
    this.padding = const EdgeInsets.symmetric(horizontal: 20),
  });

  final Widget child;
  final EdgeInsetsGeometry padding;
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: padding,
      child: Row(
        children: [
          child,
        ],
      ),
    );
  }
}

class CustomTextButton extends StatelessWidget {
  const CustomTextButton({
    super.key,
    required this.text,
    required this.onPressed,
    this.padding,
  });
  final text;
  final onPressed;
  final padding;

  @override
  Widget build(BuildContext context) {
    return IntrinsicWidth(
      child: TextButton(
        onPressed: onPressed,
        style: TextButton.styleFrom(
          padding: padding,
        ),
        child: text,
      ),
    );
  }
}

class TwoChoiceAlertDialog extends StatelessWidget {
  const TwoChoiceAlertDialog({
    super.key,
    required this.firstChoiceText,
    required this.secondChoiceText,
    required this.firstChoiceonPressed,
    required this.secondChoiceonPressed,
  });

  final Text firstChoiceText;
  final Text secondChoiceText;
  final dynamic firstChoiceonPressed;
  final dynamic secondChoiceonPressed;

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      content: const SingleChildScrollView(
        child: ListBody(
          children: <Widget>[
            Center(
              child: Text("Are you sure?", style: TextStyle(fontSize: 40, fontWeight: FontWeight.w900, color: Colors.black,)),
            )
          ],
        ),
      ),
      actions: <Widget>[
        Center(
          child: CustomButton(
            backgroundColor:const Color.fromARGB(255, 34, 34, 34),
            text: firstChoiceText,
            onPressed: firstChoiceonPressed,
            size: const Size(double.infinity, 55),
            padding: const EdgeInsets.symmetric(horizontal: 15.0),
          ),
        ),
        Center(
          child: CustomTextButton(
            text: secondChoiceText,
            onPressed: secondChoiceonPressed,
          ),
        ),
      ],
    );
  }
}

class WorkflowCard extends StatelessWidget {
  const WorkflowCard({super.key});

  @override
  Widget build(BuildContext context) {
    List<Color> colorList = [
      const Color.fromARGB(255, 46, 41, 78),
      const Color.fromARGB(255, 44, 110, 252),
      const Color.fromARGB(255, 228, 64, 95),
      const Color.fromARGB(255, 29, 185, 84),
      const Color.fromARGB(255, 205, 65, 43),
    ];

    int colorIndex = Random().nextInt(colorList.length);

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: LayoutBuilder(
        builder: (context, constraints) {
          return Container(
            decoration: BoxDecoration(
              border: Border.all(
                color: colorList[colorIndex],
                width: 2.0,
              ),
              borderRadius: BorderRadius.circular(4.0),
            ),
            child: InkWell(
              onTap: () {
                debugPrint('Card tapped.');
              },
              child: Column(
                children: [
                  SizedBox(
                    width: constraints.maxWidth,
                    height: 20,
                    child: DecoratedBox(
                      decoration: BoxDecoration(
                        color: colorList[colorIndex],
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 10),
                    child: SizedBox(
                    width: constraints.maxWidth,
                    child: const Text( "A card that can be tapped\nAnd how fit with his content",
                        style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold, color: Colors.black,),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

class ServicesCard extends StatelessWidget {
  const ServicesCard({super.key});

  @override
  Widget build(BuildContext context) {
    List<Color> colorList = [
      const Color.fromARGB(255, 41, 44, 45),
      const Color.fromARGB(255, 255, 0, 0),
      const Color.fromARGB(255, 29, 185, 84),
    ];

    List<String> nameList = [
      "Github",
      "Youtube",
      "Spotify",
    ];

    List<String> iconList = [
      "assets/icon/github.png",
      "assets/icon/youtube.png",
      "assets/icon/spotify.png",
    ];

    int serviceIndex = Random().nextInt(colorList.length);

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: LayoutBuilder(
        builder: (context, constraints) {
          return Container(
            decoration: BoxDecoration(
              color: colorList[serviceIndex],
              border: Border.all(
                color: colorList[serviceIndex],
                width: 2.0,
              ),
              borderRadius: BorderRadius.circular(4.0),
            ),

            child: InkWell(
              onTap: () {
                debugPrint('Card tapped.');
              },
              child: Column(
                children: [
                  SizedBox(
                    width: constraints.maxWidth,
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 30),
                    child: Column(
                      children: [
                        Image.asset(iconList[serviceIndex], width: 75, height: 75),
                        AutoSizeText(nameList[serviceIndex], style: const TextStyle(fontSize: 25, fontWeight: FontWeight.bold, color: Colors.white), maxLines: 1,)
                      ]
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

class TextWithBackground extends StatelessWidget {
  const TextWithBackground({
    super.key,
    required this.text,
    this.backgroundColor = Colors.white,
    this.padding = const EdgeInsets.symmetric(horizontal: 15.0),
    this.borderWidth = 1.0,
    this.borderRadius = 15.0,
    this.fontSize = 30,
    this.textColor = Colors.black,
    this.maxLines = 1,
  });

  final Color backgroundColor;
  final EdgeInsetsGeometry? padding;
  final double borderWidth;
  final double borderRadius;
  final String text;
  final double fontSize;
  final Color textColor;
  final int maxLines;

  @override
  Widget build(BuildContext context) {
    return Container (
      padding: padding,
      decoration: BoxDecoration(
        color: backgroundColor,
        border: Border.all(
          color: backgroundColor,
          width: borderWidth,
        ),
        borderRadius: BorderRadius.circular(borderRadius),
      ),
      child: AutoSizeText(
        text,
        style: TextStyle(
          fontSize: fontSize,
          fontWeight: FontWeight.bold,
          color: textColor,
          backgroundColor: backgroundColor,
        ),
        maxLines: maxLines,
      ),
    );
  }
}

class DoubleWidget extends StatelessWidget {
  const DoubleWidget({
    super.key,
    required this.firstChild,
    required this.secondChild,
    this.padding = const EdgeInsets.only(left: 50),
  });

  final Widget firstChild;
  final Widget secondChild;
  final EdgeInsetsGeometry padding;

  @override
  Widget build(BuildContext context) {
    return Row (
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        firstChild,
        Padding(padding: padding),
        secondChild,
      ],
    );
  }
}
