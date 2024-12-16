import 'package:flutter/material.dart';

import "../custom_widget.dart";
import "../footer_page.dart";
import "topBarExplore.dart";
import "exploreWidget.dart";

class ExplorePageWorkflow extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: ExplorePageWorkflowMain(),
    );
  }
}

class ExplorePageWorkflowMain extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            TopBarExplore(
              currentIndex: 1,
            ),
            MyDivider(),
            MySearchBar(),
            WorkflowCard(),
            Padding(padding: EdgeInsets.only(top: 10)),
            WorkflowCard(),
            Padding(padding: EdgeInsets.only(top: 10)),
            WorkflowCard(),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 1),
    );
  }
}
