import 'package:flutter/material.dart';

import "../custom_widget.dart";
import "../footer_page.dart";
import "topBarExplore.dart";
import "exploreWidget.dart";

class ExplorePageAll extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: ExplorePageAllMain(),
    );
  }
}

class ExplorePageAllMain extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            TopBarExplore(
              currentIndex: 0,
            ),
            MyDivider(),
            MySearchBar(),
            WorkflowCard(),
            Padding(padding: EdgeInsets.only(top: 10)),
            ServicesCard(),
            Padding(padding: EdgeInsets.only(top: 10)),
            WorkflowCard(),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 1),
    );
  }
}

class SearchBar extends StatelessWidget {
  const SearchBar({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(5),
        ),
        child: TextField(
          decoration: InputDecoration(
            prefixIcon: Padding(
              padding: EdgeInsets.only(left: 5),
              child: Icon(
                Icons.search,
                size: 35,
                color: Colors.black,
              ),
            ),
            hintStyle: TextStyle(
              fontSize: 20,
              color: const Color.fromARGB(255, 156, 156, 156),
              fontWeight: FontWeight.bold,
            ),
            hintText: 'Search',
            contentPadding: EdgeInsets.symmetric(vertical: 10, horizontal: 10),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(10),
              borderSide: BorderSide(
                color: const Color.fromARGB(255, 238, 238, 238),
                width: 5,
              ),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(10),
              borderSide: BorderSide(
                color: Colors.black,
                width: 4,
              ),
            ),
          ),
        ),
      ),
    );
  }
}
